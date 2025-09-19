import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { catchError, map, of, startWith, switchMap } from 'rxjs';

import { HousingService } from '../../services/housing.service';
import { Soba, RecenzijaSobe, Kvar, Student, MealRoomHistory } from '../../model/housing';
import { AuthService } from '../../services/auth.service';

interface ViewModel {
  loading: boolean;
  soba?: Soba;
  error?: string;
  avgOcena?: number;
  zauzeto?: number;
  slobodno?: number;
  mealHistory?: MealRoomHistory[];
}

@Component({
  selector: 'app-room-details',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './room-details.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RoomDetailsComponent {
  private route = inject(ActivatedRoute);
  private roomService = inject(HousingService);
  public auth = inject(AuthService);

  vm$ = this.route.queryParamMap.pipe(
    map(pm => pm.get('id')),
    switchMap(id => {
      if (!id) {
        return of<ViewModel>({ loading: false, error: 'Nije prosleđen ID sobe.' });
      }

      return this.roomService.getRoomDetail(id).pipe(
        switchMap(soba => {
          const zauzeto = soba.studenti?.length ?? 0;
          const slobodno = Math.max(soba.kapacitet - zauzeto, 0);

          const baseVm: ViewModel = {
            loading: false,
            soba,
            avgOcena: this.calcAvg(soba.recenzije),
            zauzeto,
            slobodno,
          };

          if (!soba.studenti?.length) {
            return of<ViewModel>({ ...baseVm, mealHistory: [] });
          }

          const usernames = soba.studenti.map(st => st.username);

          return this.roomService.getMealHistoryForUsernames(usernames).pipe(
            map(mealHistory => ({ ...baseVm, mealHistory })),
            startWith<ViewModel>({ ...baseVm, loading: true }),
            catchError(() => of<ViewModel>({ ...baseVm, error: 'Greška pri učitavanju istorije obroka.' }))
          );
        }),
        startWith<ViewModel>({ loading: true }),
        catchError(() => of<ViewModel>({ loading: false, error: 'Greška pri učitavanju sobe.' }))
      );
    })
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

  trackByMealHistory(_: number, item: MealRoomHistory) {
    return item.menu_id + item.selected_at;
  }

  canStudentAct(s: Soba | null | undefined): boolean {
    if (!s || this.auth.userRole !== 'student' || !s.studenti?.length) return false;

    const currentUsername = this.auth.username ?? '';
    return s.studenti.some(st => st.username === currentUsername);
  }

  getUniqueUsers(meals?: MealRoomHistory[]): string[] {
    if (!meals) return [];
    return Array.from(new Set(meals.map(m => m.user_name)));
  }

  getMealsForUser(userName: string, meals?: MealRoomHistory[]): MealRoomHistory[] {
    if (!meals) return [];
    return meals.filter(m => m.user_name === userName);
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }
}

