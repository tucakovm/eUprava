import {Component, OnInit, inject, ChangeDetectorRef} from '@angular/core';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';

import { User,MealHistory } from '../model/user';
import { UserService } from '../services/user.service';

@Component({
  selector: 'app-user-details',
  standalone: true,
  imports: [CommonModule, HttpClientModule, RouterModule],
  templateUrl: './user-details.html',
  styleUrls: ['./user-details.css']
})
export class UserDetails implements OnInit {
  private route = inject(ActivatedRoute);
  private userService = inject(UserService);
  private cd = inject(ChangeDetectorRef);

  user: User | null = null;
  history: MealHistory[] = [];
  loading = true;
  error: string | null = null;

  ngOnInit() {
    if (typeof window === 'undefined') {
      this.error = 'Cannot access localStorage on server';
      this.loading = false;
      return;
    }

    const storedUserId = localStorage.getItem('user');
    if (!storedUserId) {
      this.error = 'User ID not found in localStorage';
      this.loading = false;
      return;
    }

    // prvo učitavamo korisnika
    this.userService.getUserById(storedUserId).subscribe({
      next: (u) => {
        this.user = u;
        this.loading = false;
        this.cd.detectChanges();

        // posle toga učitavamo istoriju obroka
        this.loadMealHistory(storedUserId);
      },
      error: (err) => {
        this.error = 'Failed to load user';
        console.error(err);
        this.loading = false;
      }
    });
  }

  private loadMealHistory(userId: string) {
    this.userService.getMealHistory(userId).subscribe({
      next: (h) => {
        this.history = h;
        this.cd.detectChanges();
      },
      error: (err) => {
        console.error('Failed to load meal history', err);
      }
    });
  }

  deleteUser(userId: string) {
    console.log('Delete user', userId);
  }
}
