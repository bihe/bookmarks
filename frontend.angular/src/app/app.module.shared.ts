import { DragDropModule } from '@angular/cdk/drag-drop';
import { NgModule } from '@angular/core';
import { MatBadgeModule, MatCardModule, MatCheckboxModule, MatDialogModule, MatFormFieldModule, MatIconModule, MatInputModule, MatMenuModule, MatOptionModule, MatProgressSpinnerModule, MatRadioModule, MatSelectModule, MatSnackBarModule, MatTooltipModule } from '@angular/material';
import { MatButtonModule } from '@angular/material/button';
import { Title } from '@angular/platform-browser';
import { LazyLoadImageModule } from 'ng-lazyload-image';
import { AppComponent } from './components/app/app.component';
import { ConfirmDialogComponent } from './components/confirm-dialog/confirm-dialog.component';
import { DashBoardComponent } from './components/dashboard/dashboard.component';
import { FooterComponent } from './components/footer/footer.component';
import { HeaderComponent } from './components/header/header.component';
import { CreateBookmarksDialog } from './components/home/create.dialog';
import { HomeComponent } from './components/home/home.component';
import { DateFormatPipe } from './shared/pipes/dataformat';
import { EllipsisPipe } from './shared/pipes/ellipsis';
import { ApiAppInfoService } from './shared/service/api.appinfo.service';
import { ApiBookmarksService } from './shared/service/api.bookmarks.service';
import { ApplicationState } from './shared/service/application.state';


@NgModule({
  imports: [ MatProgressSpinnerModule, MatTooltipModule, MatSnackBarModule, MatButtonModule, MatDialogModule, MatInputModule, MatFormFieldModule, MatRadioModule, MatOptionModule, MatSelectModule, MatMenuModule, MatIconModule, MatBadgeModule, DragDropModule, MatCheckboxModule, MatCardModule ],
  exports: [ MatProgressSpinnerModule, MatTooltipModule, MatSnackBarModule, MatButtonModule, MatDialogModule, MatInputModule, MatFormFieldModule, MatRadioModule, MatOptionModule, MatSelectModule, MatMenuModule, MatIconModule, MatBadgeModule, DragDropModule, MatCheckboxModule, MatCardModule ],
})
export class AppMaterialModule { }

export const sharedConfig: NgModule = {
    bootstrap: [ AppComponent ],
    declarations: [
      AppComponent,
      HomeComponent,
      FooterComponent,
      HeaderComponent,
      EllipsisPipe,
      DateFormatPipe,
      CreateBookmarksDialog,
      ConfirmDialogComponent,
      DashBoardComponent
    ],
    imports: [
      AppMaterialModule,
      LazyLoadImageModule
    ],
    providers: [
      ApplicationState,
      ApiAppInfoService,
      ApiBookmarksService,
      Title
    ],
    entryComponents: [ CreateBookmarksDialog, ConfirmDialogComponent ]
};
