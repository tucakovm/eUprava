import { Routes } from '@angular/router';
import { Login } from './login/login';
import { Register } from './register/register';
import { HomeComponent } from './home/home';
import { CanteensComponent } from './canteens/canteens';
import { CanteenDetailsComponent} from './canteen-details/canteen-details';
import {MenusComponent} from './menus/menus';
import {UserDetails} from './user-details/user-details';
import { DomListComponent } from './domovi/domovi/domovi';
import { DomDetailComponent } from './domovi/dom-details/dom-details';
import { SlobodneSobeComponent } from './sobe/slobodne-sobe/slobodne-sobe';
import { DodajStudentaUSobuComponent } from './studenti/dodaj-studenta-usobu/dodaj-studenta-usobu';
import { RoomDetailsComponent } from './sobe/room-details.component/room-details.component';
import { DodajRecenzijuComponent } from './ostavi-recenziju.component/ostavi-recenziju.component';
import { PrijaviKvarComponent } from './prijavi-kvar.component/prijavi-kvar.component';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: Login },
  { path: 'register', component: Register },
  { path: 'home', component: HomeComponent },
  { path: 'canteens',component: CanteensComponent},
  { path: 'canteens/:id', component: CanteenDetailsComponent},
  { path: 'menus/:id', component: MenusComponent},
  { path: 'user-details', component: UserDetails},
  { path: 'canteens',component: CanteensComponent},
  { path: 'rooms/free', component: SlobodneSobeComponent },
  { path: 'rooms/assign', component:DodajStudentaUSobuComponent },
  {path: 'rooms/detail', component: RoomDetailsComponent},
  { path: 'rooms/review', component: DodajRecenzijuComponent },
  { path: 'rooms/fault', component: PrijaviKvarComponent },
  { path: 'doms', component: DomListComponent },
  { path: 'doms/:id', component: DomDetailComponent },

  { path: '**', redirectTo: 'login' },

];
