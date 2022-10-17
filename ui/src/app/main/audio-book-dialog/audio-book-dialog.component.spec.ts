import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AudioBookDialogComponent } from './audio-book-dialog.component';

describe('AudioBookDialogComponent', () => {
  let component: AudioBookDialogComponent;
  let fixture: ComponentFixture<AudioBookDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AudioBookDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioBookDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
