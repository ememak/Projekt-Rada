import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PollInitComponent } from './pollinit/pollinit.component';
import { VoteComponent } from './vote/vote.component';

const routes: Routes = [
  { path: 'pollinit', component: PollInitComponent },
  { path: 'vote', component: VoteComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
