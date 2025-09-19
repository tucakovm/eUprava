import { Component, OnInit } from '@angular/core';
import { CommonModule, DatePipe, CurrencyPipe } from '@angular/common';
import { take, catchError } from 'rxjs/operators';
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

  readonly today = new Date();

  constructor(private housing: HousingService) {}

  ngOnInit(): void {
    this.housing.getTodayDiningMenus()
      .pipe(
        take(1),
        catchError(err => {
          console.error('getTodayDiningMenus error', err);
          this.error = 'Greška pri učitavanju menija. Pokušaj ponovo.';
          return of<DiningMenu[]>([]);
        })
      )
      .subscribe(data => {
        this.menus = Array.isArray(data) ? data : [];
      });
  }

  trackByMenuId(_: number, m: DiningMenu) {
    return m.id;
  }
}
