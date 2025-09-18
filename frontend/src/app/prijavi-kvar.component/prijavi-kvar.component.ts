import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { HousingService } from '../services/housing.service';
import { Kvar } from '../model/housing'; 

@Component({
  selector: 'app-prijavi-kvar',
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule],
  templateUrl: './prijavi-kvar.component.html',
})
export class PrijaviKvarComponent {
  private route = inject(ActivatedRoute);
  private fb = inject(FormBuilder);
  private api = inject(HousingService);

  sobaId = this.route.snapshot.queryParamMap.get('sobaId') ?? '';
  prijavioUsername: string;
  loading = false;
  errorMsg: string | null = null;
  created?: Kvar;

  form = this.fb.group({
    opis: ['', [Validators.required, Validators.minLength(5)]],
  });

  constructor() {
    const userStr = localStorage.getItem('user');
    if (!userStr) {
      throw new Error('Nije pronađen user u localStorage!');
    }
    this.prijavioUsername = JSON.parse(userStr);
  }

  submit() {
    this.errorMsg = null;
    this.created = undefined;

    if (!this.sobaId) {
      this.errorMsg = 'Nedostaje sobaId parametar u URL-u.';
      return;
    }
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { opis } = this.form.getRawValue();
    this.loading = true;

    this.api.reportFault(this.sobaId, this.prijavioUsername, opis!).subscribe({
      next: (kvar) => {
        this.created = kvar;
        this.loading = false;
        this.form.reset(); // opcionalno: očisti formu posle uspeha
      },
      error: (err) => {
        this.loading = false;
        if (err?.status === 409) {
          this.errorMsg = err?.error || 'Već ste prijavili kvar za ovu sobu.';
        } else {
          this.errorMsg = 'Došlo je do greške pri prijavi kvara.';
        }
        console.error(err);
      },
    });
  }
}
