import { Routes } from '@angular/router';
import { Login } from './login/login';
import { Register } from './register/register';
import { HomeComponent } from './home/home';
import { CanteensComponent } from './canteens/canteens';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: Login },
  { path: 'register', component: Register },
  { path: 'home', component: HomeComponent },
  {
    path: 'canteens',component: CanteensComponent},
  { path: '**', redirectTo: 'login' }
];
