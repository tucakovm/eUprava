import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { CanteenService, CanteenDto } from '../services/canteen.service';
import {Router, RouterModule} from '@angular/router';
import { FormsModule } from '@angular/forms';


// Interfejs sa Date tipovima za frontend
interface Canteen {
  id: string;
  name: string;
  address: string;
  open_at: Date;
  close_at: Date;
}

@Component({
  selector: 'app-canteens',
  standalone: true,
  imports: [CommonModule, HttpClientModule, RouterModule, FormsModule],
  templateUrl: './canteens.html',
  styleUrls: ['./canteens.css']
})
export class CanteensComponent implements OnInit {
  private service = inject(CanteenService);
  private cd = inject(ChangeDetectorRef);
  private router = inject(Router)

  canteens: Canteen[] = [];
  loading = false;
  error: string | null = null;
  isFormOpen = false;
  newCanteen: Partial<CanteenDto> = {};

  ngOnInit(): void {
    this.getCanteens();
  }

  getCanteens() {
    this.loading = true;
    this.error = null;

    this.service.getAll().subscribe({
      next: (data: CanteenDto[]) => {
        this.canteens = data.map(c => ({
          id: c.id,
          name: c.name,
          address: c.address,
          open_at: new Date(c.open_at),
          close_at: new Date(c.close_at),
        }));
        this.loading = false;
        this.cd.detectChanges(); 
      },
      error: (err) => {
        this.error = err.message || 'Could not load canteens';
        this.loading = false;
        this.cd.detectChanges();
      }
    });
  }

  viewDetails(canteen: Canteen) {
    console.log("Selected canteen:", canteen);
    this.router.navigate(['/canteens', canteen.id]);
  }

  deleteCanteen(c: Canteen) {
    if (confirm(`Are you sure you want to delete ${c.name}?`)) {
      this.service.delete(c.id).subscribe(() => {
        this.canteens = this.canteens.filter(x => x.id !== c.id);
        this.cd.detectChanges();
      });
    }
  }

  createCanteen() {
    if (!this.newCanteen.name || !this.newCanteen.address || !this.newCanteen.open_at || !this.newCanteen.close_at) {
      alert('All fields are required.');
      return;
    }

    const openTime = this.newCanteen.open_at;
    const closeTime = this.newCanteen.close_at;

    // Validacija: open_at ne sme biti veće od close_at
    if (openTime >= closeTime) {
      alert('Open time cannot be greater than or equal to Close time.');
      return;
    }

    const payload: CanteenDto = {
      id: '',
      name: this.newCanteen.name!,
      address: this.newCanteen.address!,
      open_at: openTime,
      close_at: closeTime
    };

    this.service.create(payload).subscribe({
      next: (created) => {
        this.canteens.push({
          ...created,
          open_at: new Date(created.open_at),
          close_at: new Date(created.close_at)
        });
        this.closeForm();
      },
      error: (err) => {
        alert('Failed to create canteen: ' + (err.message || err));
      }
    });
  }


  openForm() {
    this.isFormOpen = true;
  }

  closeForm() {
    this.isFormOpen = false;
    this.newCanteen = {};
  }

}
