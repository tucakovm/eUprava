import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router'; // Dodajte ovo

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './home.html',
})
export class HomeComponent {
  private router = inject(Router);
  private auth = inject(AuthService);

  logout() {
    this.auth.logout();
    this.router.navigate(['/login']);
  }
}
