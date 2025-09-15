import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { CanteenDto, CanteenService } from '../services/canteen.service';

export interface Canteen {
  id: string;
  name: string;
  address: string;
  open_at: Date;
  close_at: Date;
}

@Component({
  selector: 'app-canteen-details',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './canteen-details.html',
  styleUrls: ['./canteen-details.css']
})
export class CanteenDetailsComponent implements OnInit {
  private service = inject(CanteenService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private cd = inject(ChangeDetectorRef);

  canteen: Canteen | null = null;
  loading = true; // Spinner startuje odmah
  error: string | null = null;

  ngOnInit(): void {
    this.loadCanteen();
  }

  loadCanteen() {
    const id = this.route.snapshot.paramMap.get('id');
    console.log('ID from route:', id);

    if (!id) {
      this.error = 'Invalid canteen ID';
      this.loading = false;
      return;
    }

    this.service.getOne(id).subscribe({
      next: (c: CanteenDto) => {
        this.canteen = {
          id: c.id,
          name: c.name,
          address: c.address,
          open_at: new Date(c.open_at),
          close_at: new Date(c.close_at)
        };
        this.loading = false;
        this.cd.detectChanges();
      },
      error: (err) => {
        this.error = err.message || 'Could not load canteen';
        this.loading = false;
        this.cd.detectChanges();
      }
    });

  }

  goBack() {
    this.router.navigate(['/canteens']);
  }

  goToMenus() {
    if (this.canteen) {
      this.router.navigate(['/menus', this.canteen.id]);
    }
  }
  
}
