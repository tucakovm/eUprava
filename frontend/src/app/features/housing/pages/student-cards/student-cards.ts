import { Component } from '@angular/core';
import { CommonModule, JsonPipe, NgIf } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators, FormGroup } from '@angular/forms';
import { HousingService } from '../../../../services/housing.service';
import { StudentskaKartica } from '../../../../model/housing';

@Component({
  selector: 'app-student-cards',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, NgIf, JsonPipe],
  templateUrl: './student-cards.html'
})
export class StudentCards {
  card?: StudentskaKartica;
  error?: string;

  createForm!: FormGroup;
  getForm!: FormGroup;
  balanceForm!: FormGroup;

  constructor(private fb: FormBuilder, private api: HousingService) {
    this.createForm = this.fb.group({
      studentId: ['', Validators.required]
    });

    this.getForm = this.fb.group({
      studentId: ['', Validators.required]
    });

    this.balanceForm = this.fb.group({
      studentId: ['', Validators.required],
      delta: [0, Validators.required]
    });
  }

  createIfMissing() { /* unchanged */ 
    if (this.createForm.invalid) return;
    this.api.createStudentCardIfMissing(this.createForm.value.studentId!).subscribe({
      next: c => this.card = c,
      error: e => this.error = e?.error || 'Error'
    });
  }

  getCard() { /* unchanged */
    if (this.getForm.invalid) return;
    this.api.getStudentCard(this.getForm.value.studentId!).subscribe({
      next: c => this.card = c,
      error: e => this.error = e?.error || 'Not found'
    });
  }

  updateBalance() { /* unchanged */
    if (this.balanceForm.invalid) return;
    const { studentId, delta } = this.balanceForm.value;
    this.api.updateStudentCardBalance(studentId!, Number(delta)).subscribe({
      next: c => this.card = c,
      error: e => this.error = e?.error || 'Error'
    });
  }
}
