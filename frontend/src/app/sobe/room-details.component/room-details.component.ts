import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { CommonModule,} from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { catchError, map, of, startWith, switchMap } from 'rxjs';

import { HousingService } from '../../services/housing.service';
import { Soba, RecenzijaSobe, Kvar, Student } from '../../model/housing';
import { AuthService } from '../../services/auth.service';

interface ViewModel {
  loading: boolean;
  soba?: Soba;
  error?: string;
  avgOcena?: number;
  zauzeto?: number;
  slobodno?: number;
}


@Component({
  selector: 'app-room-details',
  standalone: true,
  imports: [CommonModule,RouterLink],
  templateUrl: './room-details.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RoomDetailsComponent {
  private route = inject(ActivatedRoute);
  private roomService = inject(HousingService);
  public auth = inject(AuthService);

  vm$ = this.route.queryParamMap.pipe(
  map((pm) => pm.get('id')),
  switchMap((id) => {
    if (!id) {
      return of({ loading: false, error: 'Nije prosleđen ID sobe.' } as ViewModel);
    }
    return this.roomService.getRoomDetail(id).pipe(
      map((soba) => {
        const zauzeto = soba.studenti?.length ?? 0;
        const slobodno = Math.max(soba.kapacitet - zauzeto, 0);
        return {
          loading: false,
          soba,
          avgOcena: this.calcAvg(soba.recenzije),
          zauzeto,
          slobodno,
        } as ViewModel;
      }),
      startWith({ loading: true } as ViewModel),
      catchError(() =>
        of({ loading: false, error: 'Greška pri učitavanju sobe.' } as ViewModel),
      ),
    );
  }),
);


  private calcAvg(rec?: RecenzijaSobe[] | null) {
    if (!rec || rec.length === 0) return undefined;
    const sum = rec.reduce((acc, r) => acc + r.ocena, 0);
    return Math.round((sum / rec.length) * 10) / 10;
  }

  statusClass(status: Kvar['status']) {
    return {
      prijavljen: 'badge bg-warning',
      u_toku: 'badge bg-info',
      resen: 'badge bg-success',
    }[status];
  }

  trackById(_: number, item: { id: string }) {
    return item.id;
  }
  trackByStudent(_: number, s: Student) {
    return s.id;
  }
   canStudentAct(s: Soba | null | undefined): boolean {
    if (!s || this.auth.userRole !== 'student' || !s.studenti?.length) return false;

    const currentUsername = this.auth.username ?? '';

    return s.studenti.some((st: any) =>
      (typeof st?.username === 'string' && st.username === currentUsername) ||
      (typeof st?.korisnickoIme === 'string' && st.korisnickoIme === currentUsername)
    );
  }
}
