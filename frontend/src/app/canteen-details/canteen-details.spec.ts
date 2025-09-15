import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CanteenDetails } from './canteen-details';

describe('CanteenDetails', () => {
  let component: CanteenDetails;
  let fixture: ComponentFixture<CanteenDetails>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CanteenDetails]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CanteenDetails);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
