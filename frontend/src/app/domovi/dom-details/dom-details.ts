import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { HousingService } from '../../services/housing.service';
import { Dom } from '../../model/housing';
import { Observable, of } from 'rxjs';
import { catchError, map, startWith, switchMap } from 'rxjs/operators';

type DomState =
  | { status: 'loading' }
  | { status: 'success'; dom: Dom }
  | { status: 'error' };

@Component({
  selector: 'app-dom-detail',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './dom-details.html',
})
export class DomDetailComponent {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private api = inject(HousingService);

  // Jedan stream koji nosi ceo UI-state
  state$: Observable<DomState> = this.route.paramMap.pipe(
    switchMap(pm => {
      const id = pm.get('id');
      if (!id) return of<DomState>({ status: 'error' });
      return this.api.getDomById(id).pipe(
        map(dom => ({ status: 'success', dom }) as DomState),
        catchError(err => {
          console.error('GET /dom failed', err);
          return of<DomState>({ status: 'error' });
        }),
        startWith<DomState>({ status: 'loading' })
      );
    }),
    startWith<DomState>({ status: 'loading' })
  );

  back() { this.router.navigate(['/doms']); }
}
