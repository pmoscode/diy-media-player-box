import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AudioBookInformationComponent } from './audio-book-information.component';

describe('AudioBookInformationComponent', () => {
  let component: AudioBookInformationComponent;
  let fixture: ComponentFixture<AudioBookInformationComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AudioBookInformationComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioBookInformationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
