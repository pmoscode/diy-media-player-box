import {Injectable} from '@angular/core';
import {AudioBookServiceInterface} from './audio-book-service-interface';
import {AudioBook} from './audio-book';
import {Observable, of} from 'rxjs';
import {AudioBookTrackList} from './audio-book-track-list';
import {v4 as uuid} from 'uuid';
import {Card} from './card';

@Injectable({
    providedIn: 'root'
})
export class AudioBookMockService implements AudioBookServiceInterface {

    constructor() {
    }

    getAudioBooks(): Observable<AudioBook[]> {
        const audioBooks: AudioBook[] = [];

        for (let i = 0; i < 5; i++) {
            audioBooks.push(this.getRandomAudioBook(
                'AudioBook ' + i,
                6 + i,
                3 + i,
                6 + i
            ));
        }

        return of(audioBooks);
    }

    getUnassignedCardIds(): Observable<Card[]> {
        return of([{id: '9', cardId: '153245432'}, {id: '8', cardId: '13452345123'}]);
    }

    addAudioBook(audioBook: AudioBook): Observable<AudioBook> {
        return of(audioBook);
    }

    updateAudioBook(audioBook: AudioBook): Observable<AudioBook> {
        return of(audioBook);
    }

    deleteAudioBook(audioBook: AudioBook): Observable<AudioBook> {
        return of(audioBook);
    }

    deleteAllTracks(audioBook: AudioBook): Observable<AudioBook> {
        audioBook.trackList.splice(0);

        return of(audioBook);
    }

    saveTracklist(audioBook: AudioBook, files: File[]): Observable<AudioBookTrackList[]> {
        const trackList: AudioBookTrackList[] = [];

        files.forEach((file, index) => {
            const track: AudioBookTrackList = {
                title: file.name,
                track: index + 1 + audioBook.trackList.length,
                length: file.size + ''
            };

            trackList.push(track);
        });

        return of(trackList);
    }

    playTrack(audioBook: AudioBook, trackList: AudioBookTrackList): Observable<any> {
        return of({status: 'playing'});
    }

    stopTrack(): Observable<any> {
        return of({status: 'stopped'});
    }

    getRandomAudioBook(title: string, countTrackListItems: number, lastPlayedDaysBefore: number, timesPlayed: number): AudioBook {
        return {
            id: uuid(),
            card: {id: '9', cardId: '123456'},
            title,
            trackList: this.getRandomTrackList(countTrackListItems),
            lastPlayed: this.getDateBeforeDays(lastPlayedDaysBefore),
            timesPlayed
        };
    }

    getRandomTrackList(count: number): AudioBookTrackList[] {
        const trackList: AudioBookTrackList[] = [];

        for (let i = 0; i < count; i++) {
            const track: AudioBookTrackList = {
                track: i + 1,
                title: 'Track ' + (i + 1),
                length: '03:40'
            };

            trackList.push(track);
        }

        return trackList;
    }

    getDateBeforeDays(days: number) {
        const d = new Date();
        d.setDate(d.getDate() - days);

        return d;
    }
}
