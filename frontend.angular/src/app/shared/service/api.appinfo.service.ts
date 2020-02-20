import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError, timeout } from 'rxjs/operators';
import { AppInfo } from '../models/app.info.model';
import { BaseDataService } from './api.base.service';

@Injectable()
export class ApiAppInfoService extends BaseDataService {
  private readonly APP_INFO_URL: string = '/appinfo';

  constructor (private http: HttpClient) {
    super();
  }

  getApplicationInfo(): Observable<AppInfo> {
    return this.http.get<AppInfo>(this.APP_INFO_URL, this.RequestOptions)
      .pipe(
        timeout(this.RequestTimeOutDefault),
        catchError(this.handleError)
      );
  }
}
