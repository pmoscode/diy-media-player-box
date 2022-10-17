import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AudioBookTrackListComponent } from './audio-book-track-list.component';

describe('TrackListComponent', () => {
  let component: AudioBookTrackListComponent;
  let fixture: ComponentFixture<AudioBookTrackListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AudioBookTrackListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioBookTrackListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
