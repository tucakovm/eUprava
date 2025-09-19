import { Component, NgZone, ChangeDetectorRef, OnInit } from '@angular/core';
import { CommonModule, DatePipe, CurrencyPipe } from '@angular/common';
import { take, catchError, map, finalize, timeout } from 'rxjs/operators';
import { of } from 'rxjs';

import { HousingService } from '../services/housing.service';
import { DiningMenu } from '../model/housing';

@Component({
  selector: 'app-notification-meal',
  standalone: true,
  imports: [CommonModule, DatePipe, CurrencyPipe],
  templateUrl: './notification.meal.component.html',
  styleUrls: ['./notification.meal.component.css']
})
export class NotificationMealComponent implements OnInit {
  menus: DiningMenu[] = [];
  error: string | null = null;
  loading = true;

  readonly today = new Date();

  constructor(
    private housing: HousingService,
    private zone: NgZone,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.housing.getTodayDiningMenus()
      .pipe(
        timeout(5000),
        take(1),
        map((data: any) =>
          Array.isArray(data) ? data :
          Array.isArray(data?.items) ? data.items :
          Array.isArray(data?.data) ? data.data :
          Array.isArray(data?.menus) ? data.menus :
          []
        ),
        catchError(err => {
          console.error('getTodayDiningMenus error', err);
          this.error = 'Greška pri učitavanju menija. Pokušaj ponovo.';
          return of<DiningMenu[]>([]);
        }),
        finalize(() => {
          this.zone.run(() => {
            this.loading = false;
            this.cdr.markForCheck();
          });
        })
      )
      .subscribe(data => {
        this.zone.run(() => {
          this.menus = data ?? [];
          this.cdr.markForCheck();
        });
      });
  }

  trackByMenuId(_: number, m: DiningMenu) {
    return m.id;
  }
}
