import { 
  Users, 
  Briefcase, 
  ShoppingCart, 
  FileText, 
  MessageSquare, 
  Film, 
  Plane, 
  Globe, 
  CreditCard, 
  Utensils, 
  GraduationCap, 
  Trophy, 
  CheckSquare, 
  Library,
  Package,
  LucideIcon
} from 'lucide-react'

export const categoryIcons: Record<string, LucideIcon> = {
  people: Users,
  business: Briefcase,
  commerce: ShoppingCart,
  content: FileText,
  social: MessageSquare,
  media: Film,
  travel: Plane,
  location: Globe,
  finance: CreditCard,
  food: Utensils,
  education: GraduationCap,
  sports: Trophy,
  productivity: CheckSquare,
  reference: Library,
}

export function getGroupIcon(group: string) {
  return categoryIcons[group] || Package
}
