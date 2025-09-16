import { Component } from '@angular/core';
import { CommonModule, JsonPipe, NgIf } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators, FormGroup } from '@angular/forms';
import { HousingService } from '../../../../services/housing.service';
import { RecenzijaSobe } from '../../../../model/housing';

@Component({
  selector: 'app-reviews',
  standalone: true,
    imports: [CommonModule, ReactiveFormsModule, NgIf, JsonPipe],
  templateUrl: './reviews.html'
})
export class Reviews {
  review?: RecenzijaSobe;
  error?: string;

  form!: FormGroup;

  constructor(private fb: FormBuilder, private api: HousingService) {
    this.form = this.fb.group({
      sobaId: ['', Validators.required],
      autorId: ['', Validators.required],
      ocena: [5, [Validators.required, Validators.min(1), Validators.max(5)]],
      komentar: ['']
    });
  }

  submit() {
    if (this.form.invalid) return;
    const { sobaId, autorId, ocena, komentar } = this.form.value;
    this.api.addRoomReview(sobaId!, autorId!, Number(ocena), komentar || null).subscribe({
      next: r => this.review = r,
      error: e => this.error = e?.error || 'Error'
    });
  }
}
