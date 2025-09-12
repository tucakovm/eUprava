import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { LoginRequest } from '../model/auth';


@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login.html',
  styleUrl: './login.css'
})
export class Login {
  private auth = inject(AuthService);
  private router = inject(Router);

  model: LoginRequest = { identifier: '', password: '' };
  loading = false;
  error: string | null = null;

  onSubmit() {
    if (this.loading) return;
    this.error = null;
    this.loading = true;

    this.auth.login(this.model).subscribe({
      next: () => {
        this.loading = false;
        this.router.navigateByUrl('/'); // posle logina na home
      },
      error: (e) => {
        this.loading = false;
        this.error = e.message || 'Neuspešna prijava';
      }
    });
  }
}
