import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
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
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        // u storage-u je "marko123" (JSON string) -> parse vrati 'marko123'
        const parsed = JSON.parse(userStr);
        this.autorUsername =
          typeof parsed === 'string'
            ? parsed
            : parsed?.username ?? parsed?.userName ?? parsed?.email ?? '';
      } catch {
        // ako je upisan plain string bez JSON.stringify
        this.autorUsername = userStr;
      }
    }
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

    // obavezno pretvori u broj i očisti komentar
    const oc = Number(ocena);
    const kom = (komentar ?? '').trim() || null;

    this.loading = true;

    this.api
      .addRoomReview(this.sobaId, this.autorUsername, oc, kom)
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (rev) => {
          this.created = rev;
        },
        error: (err) => {
          if (err?.status === 409) {
            this.errorMsg = err?.error?.message || err?.error || 'Već ste ocenili ovu sobu.';
          } else if (err?.status === 400) {
            // backend vraća 400 za loš input/bad json
            this.errorMsg = err?.error?.message || 'Neispravan unos. Proverite podatke.';
          } else {
            this.errorMsg = 'Došlo je do greške pri ostavljanju recenzije.';
          }
          console.error(err);
        },
      });
  }
}
