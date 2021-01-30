import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PollInitComponent } from './pollinit/pollinit.component';

const routes: Routes = [
  { path: 'pollinit', component: PollInitComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
