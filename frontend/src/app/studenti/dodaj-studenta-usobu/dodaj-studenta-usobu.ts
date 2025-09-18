import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { HousingService } from '../../services/housing.service';
import { Student } from '../../model/housing';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'app-assign-student',
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule],
  templateUrl: './dodaj-studenta-usobu.html',
})
export class DodajStudentaUSobuComponent {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private fb = inject(FormBuilder);
  private api = inject(HousingService);

  domId = this.route.snapshot.queryParamMap.get('domId') ?? '';
  broj = this.route.snapshot.queryParamMap.get('broj') ?? '';

  loading = false;
  errorMsg: string | null = null;
  created?: Student;

  form = this.fb.group({
    username: ['', [Validators.required, Validators.minLength(3)]],
  });

  submit() {
    this.errorMsg = null;
    this.created = undefined;

    if (!this.domId || !this.broj) {
      this.errorMsg = 'Nedostaju parametri domId/broj u URL-u.';
      return;
    }
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { username } = this.form.getRawValue();
    this.loading = true;

    this.api.assignStudentToRoom(this.domId, this.broj, username!)
      .pipe(finalize(() => this.loading = false))
      .subscribe({
        next: (st) => {
          this.created = st;
          // odmah nazad na slobodne sobe
          this.backToFree();
        },
        error: (err) => {
          if (err?.status === 409) this.errorMsg = err?.error || 'Upis nije moguć.';
          else this.errorMsg = 'Došlo je do greške pri upisu studenta.';
          console.error(err);
        },
      });
  }

  backToFree() {
    this.router.navigate(['/rooms/free'], { queryParams: { domId: this.domId } });
  }
}
