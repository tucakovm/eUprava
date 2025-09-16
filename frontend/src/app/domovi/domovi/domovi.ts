import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule, FormControl } from '@angular/forms';
import { debounceTime, distinctUntilChanged, map, startWith, switchMap } from 'rxjs/operators';
import { Observable, combineLatest } from 'rxjs';
import { HousingService } from '../../services/housing.service';
import { Dom } from '../../model/housing';

@Component({
  selector: 'app-dom-list',
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule],
  templateUrl: './domovi.html'
})
export class DomListComponent implements OnInit {
  private api = inject(HousingService);

  search = new FormControl<string>('', { nonNullable: true });

  doms$!: Observable<Dom[]>;
  filtered$!: Observable<Dom[]>;

  ngOnInit(): void {
    this.doms$ = this.api.getAllDoms();

    this.filtered$ = combineLatest([
      this.doms$,
      this.search.valueChanges.pipe(startWith(''), debounceTime(200), distinctUntilChanged())
    ]).pipe(
      map(([doms, q]) => {
        const term = (q ?? '').toLowerCase().trim();
        if (!term) return doms;
        return doms.filter(d =>
          d.naziv.toLowerCase().includes(term) ||
          d.adresa.toLowerCase().includes(term)
        );
      })
    );
  }
}
