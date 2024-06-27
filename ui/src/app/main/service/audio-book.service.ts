import {Injectable} from '@angular/core';
import {AudioBookServiceInterface} from './audio-book-service-interface';
import {AudioBook} from './audio-book';
import {Observable} from 'rxjs';
import {AudioBookTrackList} from './audio-book-track-list';
import { HttpClient } from '@angular/common/http';
import {Card} from './card';

@Injectable({
    providedIn: 'root'
})
export class AudioBookService implements AudioBookServiceInterface {

    constructor(private httpClient: HttpClient) {
    }

    private localhost = '';

    private static checkCard(audioBook: AudioBook) {
        if (!audioBook.card) {
            audioBook.card = null;
        }
    }

    getAudioBooks(): Observable<AudioBook[]> {
        return this.httpClient.get<AudioBook[]>(this.localhost + '/api/audio-books');
    }

    getUnassignedCardIds(): Observable<Card[]> {
        return this.httpClient.get<Card[]>(this.localhost + '/api/cards/unassigned');
    }

    addAudioBook(audioBook: AudioBook): Observable<AudioBook> {
        AudioBookService.checkCard(audioBook);
        return this.httpClient.post<AudioBook>(this.localhost + '/api/audio-books', audioBook);
    }

    updateAudioBook(audioBook: AudioBook): Observable<AudioBook> {
        return this.httpClient.patch<AudioBook>(this.localhost + '/api/audio-books/' + audioBook.id, audioBook);
    }

    deleteAudioBook(audioBook: AudioBook): Observable<AudioBook> {
        return this.httpClient.delete<AudioBook>(this.localhost + '/api/audio-books/' + audioBook.id);
    }

    deleteAllTracks(audioBook: AudioBook): Observable<AudioBook> {
        return this.httpClient.delete<AudioBook>(this.localhost + '/api/audio-books/' + audioBook.id + '/tracks');
    }

    saveTracklist(audioBook: AudioBook, files: File[]): Observable<AudioBookTrackList[]> {
        if (files && files.length > 0) {
            const formData: FormData = new FormData();
            files.forEach((file, index) => {
                formData.append('file_' + index, file.slice(), file.name);
            });

            return this.httpClient.post<AudioBookTrackList[]>(this.localhost + '/api/audio-books/' + audioBook.id + '/tracks', formData);
        }
    }

    playTrack(audioBook: AudioBook, trackList: AudioBookTrackList): Observable<any> {
        return this.httpClient.post<any>(this.localhost + '/api/audio-books/' + audioBook.id + '/track/' + trackList.track + '/play', {});
    }

    stopTrack(): Observable<any> {
        return this.httpClient.post<any>(this.localhost + '/api/audio-books/stop', {});
    }
}
