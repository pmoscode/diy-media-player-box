import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AudioBookComponent } from './audio-book.component';

describe('AudioBookComponent', () => {
  let component: AudioBookComponent;
  let fixture: ComponentFixture<AudioBookComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AudioBookComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioBookComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
