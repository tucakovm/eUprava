import { Routes } from '@angular/router';
import { Login } from './login/login';
import { Register } from './register/register';
import { HomeComponent } from './home/home';
import { CanteensComponent } from './canteens/canteens';
import { CanteenDetailsComponent} from './canteen-details/canteen-details';
import {MenusComponent} from './menus/menus';
import {UserDetails} from './user-details/user-details';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: Login },
  { path: 'register', component: Register },
  { path: 'home', component: HomeComponent },
  { path: 'canteens',component: CanteensComponent},
  { path: 'canteens/:id', component: CanteenDetailsComponent},
  { path: 'menus/:id', component: MenusComponent},
  { path: 'user-details', component: UserDetails},
  { path: '**', redirectTo: 'login' }
];
