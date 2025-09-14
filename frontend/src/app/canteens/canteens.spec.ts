import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Canteens } from './canteens';

describe('Canteens', () => {
  let component: Canteens;
  let fixture: ComponentFixture<Canteens>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Canteens]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Canteens);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
