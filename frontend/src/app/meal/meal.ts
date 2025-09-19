import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MenuService, MenuWithCard } from '../services/menu.service';
import { Menu } from '../model/menus';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-meal',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './meal.html'
})
export class MealComponent implements OnInit {
  menu: Menu | null = null;
  studentCard?: { id: string; stanje: number; studentID: string };
  form!: FormGroup;
  totalPrice = 0;
  studentId = null;
  menuId = null;

  private menuService = inject(MenuService);
  private route = inject(ActivatedRoute);
  private cd = inject(ChangeDetectorRef);
  private authService = inject(AuthService);

  ngOnInit() {
     // @ts-ignore
    this.menuId = this.route.snapshot.paramMap.get('menuId');
    if (!this.menuId) {
      console.error('Menu ID not found');
      return;
    }

    // Forma
    this.form = new FormGroup({
      breakfast: new FormControl(false),
      lunch: new FormControl(false),
      dinner: new FormControl(false)
    });

    // Učitaj userId iz localStorage
    // @ts-ignore
    this.studentId = localStorage.getItem('user')
    if (!this.studentId) {
      console.error('Student ID not found in localStorage');
      return; // Ne zovi API ako nemamo studentId
    }

    this.menuService.getMenu(this.menuId, this.studentId).subscribe({
      next: (res: MenuWithCard) => {
        this.menu = res.menu;
        this.studentCard = res.card;

        this.form.reset({ breakfast: false, lunch: false, dinner: false });
        this.cd.detectChanges();
      },
      error: err => console.error('Error loading menu:', err)
    });

    this.form.valueChanges.subscribe(val => {
      this.totalPrice = 0;
      if (val.breakfast && this.menu?.breakfast) this.totalPrice += this.menu.breakfast.price;
      if (val.lunch && this.menu?.lunch) this.totalPrice += this.menu.lunch.price;
      if (val.dinner && this.menu?.dinner) this.totalPrice += this.menu.dinner.price;
    });
  }


  submit() {
    if (!this.studentCard) {
      console.error("No student card found");
      return;
    }

    if (this.studentCard.stanje < this.totalPrice) {
      alert("You do not have enough balance on your student card for this purchase!");
      return;
    }

    const payload = {
      studentUsername: this.authService.username,
      delta: -this.totalPrice,
      menuId: this.menuId,
      studentId: this.authService.userId
    };

    this.menuService.takeMeal(payload).subscribe({
      next: res => {
        console.log("Purchase successful:", res);
        alert("Meal successfully purchased!");
      },
      error: err => {
        console.error("Purchase failed:", err);
        alert("Error while purchasing meal");
      }
    });
  }

}
