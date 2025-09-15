import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { CanteenService, CanteenDto } from '../services/canteen.service';
import { RouterModule } from '@angular/router';

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
  imports: [CommonModule, HttpClientModule,RouterModule],
  templateUrl: './canteens.html',
  styleUrls: ['./canteens.css']
})
export class CanteensComponent implements OnInit {
  private service = inject(CanteenService);
  private cd = inject(ChangeDetectorRef)

  canteens: Canteen[] = [];
  loading = false;
  error: string | null = null;

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
        this.cd.detectChanges(); // forsira Angular da osveži view
      },
      error: (err) => {
        this.error = err.message || 'Could not load canteens';
        this.loading = false;
        this.cd.detectChanges();
      }
    });
  }
}
