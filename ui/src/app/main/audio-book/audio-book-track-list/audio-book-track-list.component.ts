import {Component, Inject, Input, OnInit} from '@angular/core';
import {AudioBookTrackList} from '../../service/audio-book-track-list';
import {AudioBookServiceInterface} from '../../service/audio-book-service-interface';
import {AudioBook} from '../../service/audio-book';

@Component({
    selector: 'tb-audio-book-track-list',
    templateUrl: './audio-book-track-list.component.html',
    styleUrls: ['./audio-book-track-list.component.sass']
})
export class AudioBookTrackListComponent implements OnInit {

    @Input() audioBook: AudioBook;

    trackLists: AudioBookTrackList[] = [];

    private currentlyPlaying: AudioBookTrackList;

    constructor(@Inject('AudioBookService') private audioBookService: AudioBookServiceInterface) {
    }

    ngOnInit() {
        this.trackLists = this.audioBook.trackList;
    }

    isPlaying(track: AudioBookTrackList): boolean {
        if (!this.currentlyPlaying) {
            return false;
        }

        return this.currentlyPlaying.track === track.track;
    }

    playTrack(track: AudioBookTrackList) {
        if (this.currentlyPlaying) {
            this.audioBookService.stopTrack().subscribe((status: any) => {
                this.startPlayingTrack(track);
            });
        } else {
            this.startPlayingTrack(track);
        }
    }

    startPlayingTrack(track: AudioBookTrackList) {
        this.audioBookService.playTrack(this.audioBook, track).subscribe((status: any) => {
            console.log(status);
            this.currentlyPlaying = track;
        });
    }

    stopTrack() {
        this.audioBookService.stopTrack().subscribe((status: any) => {
            console.log(status);
            this.currentlyPlaying = null;
        }, (error) => {
            console.log('Nothing was played anymore...s');
            this.currentlyPlaying = null;
        });
    }
}
