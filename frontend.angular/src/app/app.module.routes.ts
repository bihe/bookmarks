import { NgModule } from '@angular/core';
import { RouterModule, Routes, UrlMatchResult, UrlSegment } from '@angular/router';
import { DashBoardComponent } from './components/dashboard/dashboard.component';
import { HomeComponent } from './components/home/home.component';

// custom matcher
// match for all URLs starting with 'start' and collect the sub-path in
// the variable path
// e.g. /start => path: /
// e.g. /start/Folder1/Folders2 => path: /Folder1/Folder2
// e.g. /start/a/b/c/d/e/f/g => path: /a/b/c/d/e/f/g
export function matchStartAndSubPath ( url: UrlSegment[] ): UrlMatchResult {

  if (url.length === 0) {
    return null;
  }

  if (url[0].path === 'start') {
    let path = '/';
    if (url.length > 1) {
      url.forEach((e, i) => {
        if (e.path !== 'start') {
          if (!path.endsWith('/')) {
            path += '/';
          }
          path += e.path;
        }
      });
    }
    return {
      consumed: url,
      posParams: {
        path: new UrlSegment(path, {})
      }
    };
  }
  return null;
}

const routes: Routes = [
  { path: '', redirectTo: 'start', pathMatch: 'full' },
  { matcher: matchStartAndSubPath, component: HomeComponent },
  { path: 'start', component: HomeComponent },
  { path: 'dashboard', component: DashBoardComponent },
  { path: '**', redirectTo: 'start', }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
  })
export class AppRoutingModule {}
