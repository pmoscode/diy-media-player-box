import { TestBed } from '@angular/core/testing';

import { AudioBookService } from './audio-book.service';

describe('AudioBookService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: AudioBookService = TestBed.get(AudioBookService);
    expect(service).toBeTruthy();
  });
});
