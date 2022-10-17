import {AudioBookTrackList} from './audio-book-track-list';
import {Card} from './card';

export interface AudioBook {
    id: number;
    title: string;
    lastPlayed: Date;
    card?: Card;
    timesPlayed: number;
    trackList: AudioBookTrackList[];
}
