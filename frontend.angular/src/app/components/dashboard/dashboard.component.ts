import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material';
import { Title } from '@angular/platform-browser';
import { BookmarkModel } from 'src/app/shared/models/bookmarks.model';
import { ApiBookmarksService } from 'src/app/shared/service/api.bookmarks.service';
import { ApplicationState } from 'src/app/shared/service/application.state';
import { MessageUtils } from 'src/app/shared/utils/message.utils';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashBoardComponent implements OnInit {

  bookmarks: BookmarkModel[] = [];
  isUser: boolean = true;
  isAdmin: boolean = false;
  readonly MaxDashboardEntries = 45;

  constructor(private bookmarksService: ApiBookmarksService,
    private snackBar: MatSnackBar,
    private state: ApplicationState,
    private titleService: Title
  ) {}

  ngOnInit() {
    this.titleService.setTitle('bookmarks.Dashboard');
    this.state.setProgress(true);
    this.bookmarksService.getMostVisitedBookmarks(this.MaxDashboardEntries)
      .subscribe(
        data => {
          console.log(data);
          this.state.setProgress(false);
          if (data.count > 0) {
            this.bookmarks = data.value;
          } else {
            this.bookmarks = [];
          }
        },
        error => {
          this.state.setProgress(false);
          console.log('Error: ' + error);
          new MessageUtils().showError(this.snackBar, error.detail);
        }
      );

    this.state.isAdmin().subscribe(
      data => {
       this.isAdmin = data;
      }
    );
  }

  get defaultFavicon(): string {
    return 'assets/favicon.ico';
  }

  customFavicon(id: string): string {
    return `/api/v1/bookmarks/favicon/${id}`;
  }
}
