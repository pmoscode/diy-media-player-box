import {Component, Input, OnInit} from '@angular/core';
import {AudioBook} from '../../service/audio-book';

@Component({
  selector: 'tb-audio-book-information',
  templateUrl: './audio-book-information.component.html',
  styleUrls: ['./audio-book-information.component.sass']
})
export class AudioBookInformationComponent implements OnInit {

  @Input() audioBook: AudioBook;

  constructor() { }

  ngOnInit() {
  }
}
