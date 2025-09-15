import { Component, OnInit } from '@angular/core';
import { CommonModule, NgIf, NgFor } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { HousingService } from '../../../../services/housing.service';
import { Soba } from '../../../../model/housing';

@Component({
  selector: 'app-room-detail',
  standalone: true,
  imports: [CommonModule, NgIf, NgFor],
  templateUrl: './room-detail.html'
})
export class RoomDetail implements OnInit {
  room?: Soba;
  error?: string;

  constructor(private ar: ActivatedRoute, private api: HousingService) {}

  ngOnInit(): void {
    const id = this.ar.snapshot.paramMap.get('id')!;
    this.api.getRoomDetail(id).subscribe({
      next: r => this.room = r,
      error: () => this.error = 'Not found'
    });
  }
}
