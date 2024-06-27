import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {MatIconModule} from '@angular/material/icon';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MainComponent} from './main/main.component';
import {MatExpansionModule} from '@angular/material/expansion';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatTableModule} from '@angular/material/table';
import {MatButtonModule} from '@angular/material/button';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatListModule} from '@angular/material/list';
import {TimeAgoPipe} from './main/time-ago.pipe';
import {AudioBookService} from './main/service/audio-book.service';
import {environment} from '../environments/environment';
import {AudioBookTrackListComponent} from './main/audio-book/audio-book-track-list/audio-book-track-list.component';
import {AudioBookInformationComponent} from './main/audio-book/audio-book-information/audio-book-information.component';
import {AudioBookComponent} from './main/audio-book/audio-book.component';
import {AudioBookDialogComponent} from './main/audio-book-dialog/audio-book-dialog.component';
import {MatDialogModule} from '@angular/material/dialog';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatSelectModule} from '@angular/material/select';
import {MatInputModule} from '@angular/material/input';
import {MatSnackBarModule} from '@angular/material/snack-bar';
// tslint:disable-next-line:max-line-length
import {
    AudioBookManageFilesDialogComponent
} from './main/audio-book/audio-book-manage-files-dialog/audio-book-manage-files-dialog.component';
import {NgxFileDropModule} from 'ngx-file-drop';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatCheckboxModule} from '@angular/material/checkbox';

@NgModule({ declarations: [
        AppComponent,
        MainComponent,
        TimeAgoPipe,
        AudioBookTrackListComponent,
        AudioBookInformationComponent,
        AudioBookComponent,
        AudioBookDialogComponent,
        AudioBookManageFilesDialogComponent
    ],
    bootstrap: [AppComponent], imports: [BrowserModule,
        AppRoutingModule,
        BrowserAnimationsModule,
        FormsModule,
        MatIconModule,
        MatExpansionModule,
        MatToolbarModule,
        MatTableModule,
        MatButtonModule,
        MatGridListModule,
        MatListModule,
        MatDialogModule,
        MatFormFieldModule,
        MatSelectModule,
        ReactiveFormsModule,
        MatInputModule,
        MatSnackBarModule,
        NgxFileDropModule,
        MatProgressSpinnerModule,
        MatCheckboxModule], providers: [
        { provide: 'AudioBookService', useClass: environment.audioBookService },
        provideHttpClient(withInterceptorsFromDi())
    ] })
export class AppModule {
}
