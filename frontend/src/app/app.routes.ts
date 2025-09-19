import { Routes } from '@angular/router';
import { Login } from './login/login';
import { Register } from './register/register';
import { HomeComponent } from './home/home';
import { CanteensComponent } from './canteens/canteens';
import { CanteenDetailsComponent } from './canteen-details/canteen-details';
import { MenusComponent } from './menus/menus';
import { UserDetails } from './user-details/user-details';
import { MealComponent } from './meal/meal';
import { DomListComponent } from './domovi/domovi/domovi';
import { DomDetailComponent } from './domovi/dom-details/dom-details';
import { SlobodneSobeComponent } from './sobe/slobodne-sobe/slobodne-sobe';
import { DodajStudentaUSobuComponent } from './studenti/dodaj-studenta-usobu/dodaj-studenta-usobu';
import { RoomDetailsComponent } from './sobe/room-details.component/room-details.component';
import { DodajRecenzijuComponent } from './ostavi-recenziju.component/ostavi-recenziju.component';
import { PrijaviKvarComponent } from './prijavi-kvar.component/prijavi-kvar.component';
import { UnauthorizedComponent } from './unauthorized/unauthorized.component';

// ⬇️ uvezi tvoj guard
import { roleCanActivate, roleCanMatch } from './auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },

  // javne rute
  { path: 'login', component: Login },
  { path: 'register', component: Register },

  // zajedničke rute (admin + student)
  {
    path: 'home',
    component: HomeComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'canteens',
    component: CanteensComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'canteens/:id',
    component: CanteenDetailsComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'menus/:id',
    component: MenusComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'user-details',
    component: UserDetails,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'doms',
    component: DomListComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'doms/:id',
    component: DomDetailComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },
  {
    path: 'rooms/detail',
    component: RoomDetailsComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin', 'student'] }
  },

  // samo ADMIN
  {
    path: 'rooms/free',
    component: SlobodneSobeComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin','student'] }
  },
  {
    path: 'rooms/assign',
    component: DodajStudentaUSobuComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['admin'] }
  },

  // samo STUDENT
  {
    path: 'meal/:menuId',
    component: MealComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['student'] }
  },
  {
    path: 'rooms/review',
    component: DodajRecenzijuComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['student'] }
  },
  {
    path: 'rooms/fault',
    component: PrijaviKvarComponent,
    canMatch: [roleCanMatch],
    canActivate: [roleCanActivate],
    data: { roles: ['student'] }
  },

  // 401/403 stranica
  { path: 'unauthorized', component: UnauthorizedComponent },

  // fallback
  { path: '**', redirectTo: 'login' },
];
