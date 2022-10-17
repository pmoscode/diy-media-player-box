import {Component, EventEmitter, Inject, Input, OnInit, Output} from '@angular/core';
import {AudioBook} from '../service/audio-book';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {AudioBookServiceInterface} from '../service/audio-book-service-interface';
import {MatSnackBar} from '@angular/material/snack-bar';
import {AudioBookDialogComponent} from '../audio-book-dialog/audio-book-dialog.component';
import {AudioBookManageFilesDialogComponent} from './audio-book-manage-files-dialog/audio-book-manage-files-dialog.component';
import {AudioBookTrackList} from '../service/audio-book-track-list';
import {Subject} from 'rxjs';
import {indicate} from '../../observable-helper';

@Component({
    selector: 'tb-audio-book',
    templateUrl: './audio-book.component.html',
    styleUrls: ['./audio-book.component.sass']
})
export class AudioBookComponent implements OnInit {

    @Input() audioBook: AudioBook;

    @Output() audioBookDeletedEvent = new EventEmitter();

    loadingTracks$ = new Subject<boolean>();

    constructor(private dialog: MatDialog,
                @Inject('AudioBookService') private audioBookService: AudioBookServiceInterface,
                private snackBar: MatSnackBar) {
    }

    ngOnInit() {
    }

    editAudioBook(audioBook: AudioBook) {
        const dialogRef: MatDialogRef<AudioBookDialogComponent> = this.dialog.open(AudioBookDialogComponent, {
            width: '600px',
            data: {audioBook}
        });

        dialogRef.afterClosed().subscribe((data: { audioBook: AudioBook }) => {
            if (data && data.audioBook) {
                this.updateAudioBook(data.audioBook);
            }
        });
    }

    private updateAudioBook(audioBook: AudioBook) {
        const mergedAudioBook = Object.assign(this.audioBook, audioBook);
        this.audioBookService.updateAudioBook(mergedAudioBook).subscribe((updatedAudioBook: AudioBook) =>
                this.snackBar.open('Audio book saved: ' + updatedAudioBook.title, '', {duration: 5000}),
            (error: Error) => {
                console.log(error);
                this.snackBar.open('Audio book not saved...!', '', {duration: 5000});
            });
    }

    deleteAudioBook(audioBook: AudioBook) {
        this.audioBookService.deleteAudioBook(audioBook).subscribe((deletedAudioBook: AudioBook) => {
                this.snackBar.open('Audio book deleted: ' + deletedAudioBook.title, '', {duration: 5000});
                this.audioBookDeletedEvent.emit();
            },
            (error: Error) => {
                console.log(error);
                this.snackBar.open('Audio book not deleted...!', '', {duration: 5000});
            });
    }

    manageAudioFiles(audioBook: AudioBook) {
        const hasFiles: boolean = audioBook.trackList ? audioBook.trackList.length > 0 : false;
        const dialogRef: MatDialogRef<AudioBookManageFilesDialogComponent> = this.dialog.open(AudioBookManageFilesDialogComponent, {
            width: '800px',
            data: {hasFiles}
        });

        dialogRef.componentInstance.deleteAllTracksEvent.subscribe(() => {
            this.audioBookService.deleteAllTracks(audioBook).subscribe((updatedAudioBook: AudioBook) => {
                    this.snackBar.open('All tracks from audio book deleted', '', {duration: 5000});
                    this.audioBook = updatedAudioBook;
                },
                (error: Error) => {
                    console.log(error);
                    this.snackBar.open('No tracks from audio book deleted...!', '', {duration: 5000});
                });
        });

        dialogRef.afterClosed().subscribe((data: { files: File[] }) => {
            if (data && data.files) {
                this.audioBookService.saveTracklist(audioBook, data.files)
                    .pipe(indicate(this.loadingTracks$))
                    .subscribe((trackList: AudioBookTrackList[]) => {
                        this.snackBar.open('All tracks from audio book saved', '', {duration: 5000});
                        trackList.forEach(track => console.log('Track: ', track));
                        trackList.forEach(track => audioBook.trackList.push(track));
                        audioBook.trackList = audioBook.trackList.slice(); // Notifies Angular on changed contents
                    },
                    (error: Error) => {
                        console.log(error);
                        this.snackBar.open('No tracks from audio book saved...!', '', {duration: 5000});
                    });
            }
        });
    }
}
