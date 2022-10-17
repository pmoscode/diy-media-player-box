import {AudioBook} from './audio-book';
import {Observable} from 'rxjs';
import {AudioBookTrackList} from './audio-book-track-list';
import {Card} from './card';

export interface AudioBookServiceInterface {
    getAudioBooks(): Observable<AudioBook[]>;
    getUnassignedCardIds(): Observable<Card[]>;
    addAudioBook(audioBook: AudioBook): Observable<AudioBook>;
    updateAudioBook(audioBook: AudioBook): Observable<AudioBook>;
    saveTracklist(audioBook: AudioBook, files: File[]): Observable<AudioBookTrackList[]>;
    deleteAudioBook(audioBook: AudioBook): Observable<AudioBook>;
    deleteAllTracks(audioBook: AudioBook): Observable<AudioBook>;
    playTrack(audioBook: AudioBook, trackList: AudioBookTrackList): Observable<any>;
    stopTrack(): Observable<any>;
}
