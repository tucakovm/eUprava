import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router'; // ⬅️ dodaj Router
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { HousingService } from '../services/housing.service';
import { RecenzijaSobe } from '../model/housing';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'app-dodaj-recenziju',
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule],
  templateUrl: './ostavi-recenziju.component.html',
})
export class DodajRecenzijuComponent {
  private route = inject(ActivatedRoute);
  private fb = inject(FormBuilder);
  private api = inject(HousingService);
  private router = inject(Router); // ⬅️ inject Router

  sobaId = this.route.snapshot.queryParamMap.get('sobaId') ?? '';
  autorUsername = '';
  loading = false;
  errorMsg: string | null = null;
  created?: RecenzijaSobe;

  form = this.fb.group({
    ocena: this.fb.control<number | null>(null, [
      Validators.required,
      Validators.min(1),
      Validators.max(5),
    ]),
    komentar: this.fb.control<string>(''),
  });

  constructor() {
    const userStr = localStorage.getItem('user') ?? '';
    let autorUsername = '';

    if (userStr) {
      const trimmed = userStr.trim();
      const looksLikeJson = trimmed.startsWith('{') || trimmed.startsWith('[') || trimmed.startsWith('"');

      if (looksLikeJson) {
        try {
          const parsed = JSON.parse(userStr);
          autorUsername =
            typeof parsed === 'string'
              ? parsed
              : parsed?.username ?? parsed?.userName ?? parsed?.email ?? '';
        } catch {
          autorUsername = userStr;
        }
      } else {
        autorUsername = userStr;
      }
    }

    this.autorUsername = autorUsername;

    if (!this.autorUsername) {
      this.errorMsg = 'Niste prijavljeni. Ulogujte se pa pokušajte ponovo.';
    }
  }

  submit() {
    this.errorMsg = null;
    this.created = undefined;

    if (!this.sobaId) {
      this.errorMsg = 'Nedostaje sobaId parametar u URL-u.';
      return;
    }
    if (!this.autorUsername) {
      this.errorMsg = 'Nedostaje korisničko ime.';
      return;
    }
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { ocena, komentar } = this.form.getRawValue();
    const oc = Number(ocena);
    const kom = (komentar ?? '').trim() || null;

    this.loading = true;

    this.api
      .addRoomReview(this.sobaId, this.autorUsername, oc, kom)
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (rev) => {
          const targetId = rev?.sobaId ?? this.sobaId;

          this.router.navigateByUrl(`/rooms/detail?id=${encodeURIComponent(targetId)}`);
        },
        error: (err) => {
          if (err?.status === 409) {
            this.errorMsg = err?.error?.message || err?.error || 'Već ste ocenili ovu sobu.';
          } else if (err?.status === 400) {
            this.errorMsg = err?.error?.message || 'Neispravan unos. Proverite podatke.';
          } else {
            this.errorMsg = 'Došlo je do greške pri ostavljanju recenzije.';
          }
          console.error(err);
        },
      });
  }
}
