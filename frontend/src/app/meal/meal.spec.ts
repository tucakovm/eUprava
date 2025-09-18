import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Meal } from './meal';

describe('Meal', () => {
  let component: Meal;
  let fixture: ComponentFixture<Meal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Meal]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Meal);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
