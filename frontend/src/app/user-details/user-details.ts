import {Component, OnInit, inject, ChangeDetectorRef} from '@angular/core';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

import { User, MealHistory, MenuReview } from '../model/user';
import { UserService } from '../services/user.service';
import {AuthService} from '../services/auth.service';
import {StudentskaKartica} from '../model/housing';
import {HousingService} from '../services/housing.service';

declare var bootstrap: any;

@Component({
  selector: 'app-user-details',
  standalone: true,
  imports: [CommonModule, HttpClientModule, RouterModule, FormsModule],
  templateUrl: './user-details.html',
  styleUrls: ['./user-details.css']
})
export class UserDetails implements OnInit {
  private route = inject(ActivatedRoute);
  private userService = inject(UserService);
  private cd = inject(ChangeDetectorRef);
  private authService = inject(AuthService)
  private housingService = inject(HousingService);


  user: User | null = null;
  history: MealHistory[] = [];
  loading = true;
  error: string | null = null;
  isAdmin: boolean = false;
  creatingCard = false;
  cardError: string | null = null;
  createdCard?: StudentskaKartica;


  // Review modal data
  selectedMeal: MealHistory | null = null;
  reviewData: MenuReview = {
    id: '',
    menu_id: '',
    user_id: '',
    breakfast_review: 0,
    lunch_review: 0,
    dinner_review: 0
  };
  reviewLoading = false;
  private reviewModal: any;

  ngOnInit() {
    this.isAdmin = this.authService.userRole === 'admin';
    if (typeof window === 'undefined') {
      this.error = 'Cannot access localStorage on server';
      this.loading = false;
      return;
    }

    const storedUserId = localStorage.getItem('user');
    const loadMealHistoryparam = localStorage.getItem('userId')
    if (!storedUserId) {
      this.error = 'User ID not found in localStorage';
      this.loading = false;
      return;
    }

    // Load user first
    this.userService.getUserById(storedUserId,loadMealHistoryparam).subscribe({
      next: (u) => {
        this.user = u;
        this.loading = false;
        this.cd.detectChanges();

        // Then load meal history with reviews
        this.loadMealHistory(loadMealHistoryparam);
      },
      error: (err) => {
        this.error = 'Failed to load user';
        console.error(err);
        this.loading = false;
      }
    });

    // Initialize Bootstrap modal after view init
    setTimeout(() => {
      if (typeof bootstrap !== 'undefined') {
        const modalElement = document.getElementById('reviewModal');
        if (modalElement) {
          this.reviewModal = new bootstrap.Modal(modalElement);
        }
      }
    }, 100);
  }

  private loadMealHistory(userId: string | null) {
    this.userService.getMealHistoryWithReviews(userId).subscribe({
      next: (h) => {
        this.history = h;
        this.cd.detectChanges();
      },
      error: (err) => {
        console.error('Failed to load meal history');
      }
    });
  }

  openReviewModal(meal: MealHistory) {
    this.selectedMeal = meal;

    if (meal.review) {
      // Edit existing review
      this.reviewData = { ...meal.review };
    } else {
      // Create new review
      this.reviewData = {
        id: '',
        menu_id: meal.menu_id || '',
        user_id: this.user?.id || '',
        breakfast_review: 0,
        lunch_review: 0,
        dinner_review: 0
      };
    }
  }

  setRating(field: keyof MenuReview, rating: number) {
    if (field === 'breakfast_review' || field === 'lunch_review' || field === 'dinner_review') {
      this.reviewData[field] = rating;
    }
  }

  getStarRating(rating: number): string {
    return '★'.repeat(rating) + '☆'.repeat(5 - rating);
  }

  saveReview() {
    if (!this.selectedMeal || !this.user) return;

    this.reviewLoading = true;

    const isUpdate = !!this.selectedMeal.review;

    const reviewObservable = isUpdate
      ? this.userService.updateReview(this.reviewData)
      : this.userService.createReview(this.reviewData);

    reviewObservable.subscribe({
      next: (review) => {
        // Update the meal's review in the history
        if (this.selectedMeal) {
          this.selectedMeal.review = review;
        }

        this.reviewLoading = false;
        this.closeModal();
        this.cd.detectChanges();
      },
      error: (err) => {
        console.error('Failed to save review', err);
        this.reviewLoading = false;
        // You might want to show an error message here
      }
    });
  }

  private closeModal() {
    if (this.reviewModal) {
      this.reviewModal.hide();
    }
  }

  deleteUser(userId: string) {
    console.log('Delete user', userId);
  }

  // === NEW: create student card ===
  createCard() {
    this.cardError = null;
    this.createdCard = undefined;

    if (!this.user?.username) {
      this.cardError = 'Username is missing for this user.';
      return;
    }

    this.creatingCard = true;
    this.housingService.createStudentCardIfMissing(this.user.username).subscribe({
      next: (card) => {
        this.createdCard = card;
        this.creatingCard = false;
        this.cd.detectChanges();
      },
      error: (err) => {
        this.creatingCard = false;
        this.cardError =
          err?.error?.message || err?.message || 'Failed to create student card.';
        console.error('createStudentCardIfMissing error', err);
      }
    });
  }
}
