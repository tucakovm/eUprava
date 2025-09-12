import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { RegisterRequest } from '../model/auth';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './register.html',
  styleUrl: './register.css'
})
export class Register {
  private auth = inject(AuthService);
  private router = inject(Router);

  model: RegisterRequest = {
    firstname: '',
    lastname: '',
    username: '',
    email: '',
    password: ''
  };

  loading = false;
  error: string | null = null;
  success = false;

  onSubmit() {
    if (this.loading) return;
    this.error = null;
    this.success = false;
    this.loading = true;

    this.auth.register(this.model).subscribe({
      next: () => {
        this.loading = false;
        this.success = true;
        // opcionalno: automatski na login
        setTimeout(() => this.router.navigateByUrl('/login'), 1000);
      },
      error: (e) => {
        this.loading = false;
        this.error = e.message || 'Registracija neuspešna';
      }
    });
  }
}
