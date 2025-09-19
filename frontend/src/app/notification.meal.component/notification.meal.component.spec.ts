import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NotificationMealComponent } from './notification.meal.component';

describe('NotificationMealComponent', () => {
  let component: NotificationMealComponent;
  let fixture: ComponentFixture<NotificationMealComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NotificationMealComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NotificationMealComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
