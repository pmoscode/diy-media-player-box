import {Component, Inject, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {AudioBookServiceInterface} from '../service/audio-book-service-interface';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {AudioBook} from '../service/audio-book';
import {Card} from '../service/card';

@Component({
    selector: 'tb-audio-book-dialog',
    templateUrl: './audio-book-dialog.component.html',
    styleUrls: ['./audio-book-dialog.component.sass']
})
export class AudioBookDialogComponent implements OnInit {

    formGroup = new FormGroup({
        title: new FormControl('', [Validators.required]),
        card: new FormControl('', [])
    });

    cards: Card[] = [];

    constructor(@Inject('AudioBookService') private audioBookService: AudioBookServiceInterface,
                public dialogRef: MatDialogRef<AudioBookDialogComponent>,
                @Inject(MAT_DIALOG_DATA) public data: { audioBook: AudioBook }) {
    }

    ngOnInit() {
        this.audioBookService.getUnassignedCardIds().subscribe((cards: Card[]) => {
            this.cards = cards;
            if (this.data.audioBook) {
                if (this.data.audioBook.card) {
                    this.cards.push(this.data.audioBook.card);
                }
                this.formGroup.controls.title.setValue(this.data.audioBook.title);
                this.formGroup.controls.card.setValue(this.data.audioBook.card);

                this.cards.sort((a, b) => a.cardId.localeCompare(b.cardId));
            }
        });
    }

    compareCard(card1: Card, card2: Card): boolean {
        return card1.id === card2.id;
    }

    getTitleMessage() {
        if (this.formGroup.controls.title.hasError('required')) {
            return 'Audio book title is required';
        }
        return '';
    }

    cancel() {
        this.dialogRef.close();
    }
}
