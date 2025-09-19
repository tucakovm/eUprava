import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { HousingService } from '../../services/housing.service';
import { Soba } from '../../model/housing';
import { Observable, of } from 'rxjs';
import { catchError, map, startWith, switchMap } from 'rxjs/operators';
import { AuthService } from '../../services/auth.service';

type RoomsState =
  | { status: 'loading'; domId: string | null }
  | { status: 'success'; domId: string; rooms: Soba[] }
  | { status: 'error'; domId: string | null }
  | { status: 'missing' };

@Component({
  selector: 'app-free-rooms',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './slobodne-sobe.html',
})
export class SlobodneSobeComponent {
  private route = inject(ActivatedRoute);
  private api = inject(HousingService);
  private router = inject(Router);
  public auth = inject(AuthService);

  state$: Observable<RoomsState> = this.route.queryParamMap.pipe(
    switchMap(qp => {
      const domId = qp.get('domId');
      if (!domId) return of<RoomsState>({ status: 'missing' });
      return this.api.listFreeRooms(domId).pipe(
        map(rooms => ({ status: 'success', domId, rooms }) as RoomsState),
        catchError(err => {
          console.error('GET /rooms/free failed', err);
          return of<RoomsState>({ status: 'error', domId });
        }),
        startWith<RoomsState>({ status: 'loading', domId })
      );
    }),
    startWith<RoomsState>({ status: 'loading', domId: null })
  );

  toAssign(domId: string, broj: string) {
    this.router.navigate(['/rooms/assign'], { queryParams: { domId, broj } });
  }
  toDetails(id:string) {
    this.router.navigate(['/rooms/detail'], { queryParams: { id } });
  }

}
