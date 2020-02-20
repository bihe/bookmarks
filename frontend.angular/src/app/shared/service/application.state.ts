import { ReplaySubject } from 'rxjs';
import { AppInfo } from '../models/app.info.model';

export class ApplicationState {
    private progress: ReplaySubject<boolean> = new ReplaySubject();
    private appInfo: ReplaySubject<AppInfo> = new ReplaySubject();
    private admin: ReplaySubject<boolean> = new ReplaySubject();

    public setAppInfo(data: AppInfo) {
        this.appInfo.next(data);
    }

    public getAppInfo(): ReplaySubject<AppInfo> {
        return this.appInfo;
    }

    public setProgress(data: boolean) {
        this.progress.next(data);
    }

    public getProgress(): ReplaySubject<boolean> {
        return this.progress;
    }

    public isAdmin(): ReplaySubject<boolean> {
      return this.admin;
    }

    public setAdmin(data: boolean) {
      this.admin.next(data);
    }
}
