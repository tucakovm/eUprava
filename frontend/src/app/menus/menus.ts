import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule, Router, ActivatedRoute } from '@angular/router';
import { Menu, Weekday } from '../model/menus';
import { MenuService } from '../services/menu.service';

@Component({
  selector: 'app-menus',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './menus.html',
  styleUrls: ['./menus.css']
})
export class MenusComponent implements OnInit {
  private menuService = inject(MenuService);
  private cd = inject(ChangeDetectorRef);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  menus: Menu[] = [];
  loading = false;
  error: string | null = null;
  canteenId: string | null = null;

  // Modal kontrola
  isMenuFormOpen = false;

  newMenu: Menu = {
    id: '',
    name: '',
    weekday: Weekday.Monday,
    canteen_id: this.canteenId ?? '',
    breakfast: { name: '', description: '', price: 0 },
    lunch: { name: '', description: '', price: 0 },
    dinner: { name: '', description: '', price: 0 }
  };


  weekdays = Object.values(Weekday);

  ngOnInit(): void {
    this.canteenId = this.route.snapshot.paramMap.get('id');
    if (this.canteenId) {
      this.getMenus(this.canteenId);
      this.newMenu.canteen_id = this.canteenId;
    }
  }


  getMenus(canteenId: string) {
    this.loading = true;
    this.error = null;
    this.menuService.getAll(canteenId).subscribe({
      next: (data) => {
        this.menus = data;
        this.loading = false;
        this.cd.detectChanges();
      },
      error: (err) => {
        this.error = err.message || 'Failed to load menus';
        this.loading = false;
      }
    });
  }

  openMenuForm() {
    this.isMenuFormOpen = true;
  }

  closeMenuForm() {
    this.isMenuFormOpen = false;
    this.newMenu = {
      id: '',
      name: '',
      weekday: Weekday.Monday,
      canteen_id: this.canteenId ?? '',
      breakfast: { name: '', description: '', price: 0 },
      lunch: { name: '', description: '', price: 0 },
      dinner: { name: '', description: '', price: 0 }
    };
  }


  createMenu() {
    if (!this.newMenu.name || !this.newMenu.weekday) {
      alert('Menu name and weekday are required.');
      return;
    }

    if (!this.newMenu.breakfast?.name || !this.newMenu.lunch?.name || !this.newMenu.dinner?.name) {
      alert('All meals must have a name.');
      return;
    }

    this.menuService.create(this.newMenu as Menu).subscribe({
      next: (created) => {
        console.log("menu push", created)
        this.menus.push(created);
        this.closeMenuForm();
      },
      error: (err) => {
        alert('Failed to create menu: ' + (err.message || err));
      }
    });
  }

  viewMenu(menu: Menu) {
    this.router.navigate(['/menus', menu.id]);
  }

  deleteMenu(menu: Menu) {
    // if (confirm(`Are you sure you want to delete menu "${menu.name}"?`)) {
    //   this.menuService.delete(menu.id).subscribe(() => {
    //     this.menus = this.menus.filter(m => m.id !== menu.id);
    //     this.cd.detectChanges();
    //   });
    // }
  }
}
