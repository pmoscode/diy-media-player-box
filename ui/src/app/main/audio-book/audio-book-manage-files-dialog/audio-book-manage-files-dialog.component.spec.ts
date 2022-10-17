import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AudioBookManageFilesDialogComponent } from './audio-book-manage-files-dialog.component';

describe('AudioBookManageFilesDialogComponent', () => {
  let component: AudioBookManageFilesDialogComponent;
  let fixture: ComponentFixture<AudioBookManageFilesDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AudioBookManageFilesDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioBookManageFilesDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
