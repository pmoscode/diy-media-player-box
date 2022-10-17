import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {FileSystemFileEntry, NgxFileDropEntry} from 'ngx-file-drop';

@Component({
    selector: 'tb-audio-book-manage-files-dialog',
    templateUrl: './audio-book-manage-files-dialog.component.html',
    styleUrls: ['./audio-book-manage-files-dialog.component.sass']
})
export class AudioBookManageFilesDialogComponent implements OnInit {

    @Output() deleteAllTracksEvent = new EventEmitter();

    files: File[] = [];

    constructor(public dialogRef: MatDialogRef<AudioBookManageFilesDialogComponent>,
                @Inject(MAT_DIALOG_DATA) public data: { hasFiles: boolean }) {
    }

    ngOnInit() {
    }

    public dropped(files: NgxFileDropEntry[]) {
        for (const droppedFile of files) {
            if (droppedFile.fileEntry.isFile) {
                const fileEntry = droppedFile.fileEntry as FileSystemFileEntry;
                fileEntry.file((file: File) => {
                    if (file.type === 'audio/mpeg') {
                        this.files.push(file);
                    } else {
                        console.log(file.type, ' is not mp3 file');
                    }
                });
            } else {
                console.log(droppedFile.relativePath, 'no directories allowed');
            }
        }
    }

    close() {
        this.dialogRef.close();
    }

    deleteAllTracks() {
        this.deleteAllTracksEvent.emit();
    }
}
