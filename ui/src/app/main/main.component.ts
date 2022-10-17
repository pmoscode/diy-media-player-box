import {Component, Inject, OnInit} from '@angular/core';
import {AudioBookServiceInterface} from './service/audio-book-service-interface';
import {AudioBook} from './service/audio-book';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {AudioBookDialogComponent} from './audio-book-dialog/audio-book-dialog.component';
import {MatSnackBar} from '@angular/material/snack-bar';

@Component({
    selector: 'tb-main',
    templateUrl: './main.component.html',
    styleUrls: ['./main.component.sass']
})
export class MainComponent implements OnInit {

    audioBooks: AudioBook[] = [];

    constructor(private dialog: MatDialog,
                @Inject('AudioBookService') private audioBookService: AudioBookServiceInterface,
                private snackBar: MatSnackBar) {
    }

    ngOnInit() {
        this.loadAudioBooks();
    }

    addNewAudioBookDialog() {
        const dialogRef: MatDialogRef<AudioBookDialogComponent> = this.dialog.open(AudioBookDialogComponent, {
            width: '600px',
            data: {audioBook: undefined}
        });

        dialogRef.afterClosed().subscribe((data: { audioBook: AudioBook }) => {
            if (data && data.audioBook) {
                this.addNewAudioBook(data.audioBook);
            }
        });
    }

    private addNewAudioBook(audioBook: AudioBook) {
        this.audioBookService.addAudioBook(audioBook).subscribe((addedAudioBook: AudioBook) => {
                this.loadAudioBooks();
                this.snackBar.open('Audio book saved: ' + addedAudioBook.title, '', {duration: 5000});
            },
            (error: Error) => {
                console.log(error);
                this.snackBar.open('Audio book not saved...!', '', {duration: 5000});
            });
    }

    audioBookDeleted() {
        this.loadAudioBooks()
    }

    private loadAudioBooks() {
        this.audioBookService.getAudioBooks().subscribe((audioBooks: AudioBook[]) => {
            this.audioBooks = audioBooks;
        });
    }
}
